let AWS = require('aws-sdk');

let asCient = new AWS.AutoScaling();
let ecsClient = new AWS.ECS();

let eventProcessor = (event, context) => {
    event.Records.forEach(recordProcessor);
};
let recordProcessor = record => {
    console.log(`Record: ${JSON.stringify(record)}`);
    var snsMessage = record.Sns.Message;
    var lifeCycleActionInfo = JSON.parse(snsMessage);

    var params = {
        cluster: process.env.CLUSTER_NAME,
        status: 'ACTIVE'
    };
    ecsClient.listContainerInstances(params, function(err, data) {
        if (err) {
            console.error(`List Instance Failed: ${JSON.stringify(err)}`);
            return;
        }
        console.log(`Cluster Instances: ${JSON.stringify(data)}`);
        var params = {
            containerInstances: data.containerInstanceArns,
            cluster: process.env.CLUSTER_NAME
        };
        ecsClient.describeContainerInstances(params, (describeError, describeData) => {
            if (describeError) {
                console.error(`Describe Instances Failed: ${JSON.stringify(describeError)}`);
                return;
            }

            describeData.containerInstances.forEach(instance => {
                if (instance.ec2InstanceId !== lifeCycleActionInfo.EC2InstanceId)
                    return;

                drainInstance(instance.containerInstanceArn)
                    .then(() => waitInstanceDrain(instance.containerInstanceArn))
                    .then(() => completeLca(lifeCycleActionInfo));
            });
        });
    });

};
let drainInstance = (cluster => {
    return instanceId => {
        return new Promise((fulfill, reject) => {
            var params = {
                containerInstances: [instanceId],
                status: 'DRAINING',
                cluster: cluster
            };
            console.log(`Draing Instance Param: ${JSON.stringify(params)}`);
            ecsClient.updateContainerInstancesState(params, (err, data) => {
                if (err) reject(err);
                else fulfill(data);
            });
        });
    };
})(process.env.CLUSTER_NAME);

let waitInstanceDrain = (cluster => {
    return instanceId => {
        return new Promise((fulfill, reject) => {
            var params = {
                containerInstances: [instanceId],
                cluster: cluster
            };
            console.log(`Wait Instance Param: ${JSON.stringify(params)}`);
            var count = 0;
            let describeContainerInstanceCallback = (err, data) => {
                if (err) reject(err);
                else if (data.containerInstances[0].runningTasksCount === 0) {
                    fulfill(data);
                // Lambda max timeout = 300 sec
                } else if (count >= 10) {
                    fulfill(data);
                } else {
                    // wait 15 seconds and check again
                    count++;
                    console.log(`Waiting... Attempt: ${count}`);
                    setTimeout(() => ecsClient.describeContainerInstances(params, describeContainerInstanceCallback), 29000);
                }
            };
            ecsClient.describeContainerInstances(params, describeContainerInstanceCallback);
        });
    };
})(process.env.CLUSTER_NAME);

let completeLca = lifeCycleActionInfo => {
    return new Promise((fulfill, reject) => {
        var params = {
            AutoScalingGroupName: lifeCycleActionInfo.AutoScalingGroupName,
            LifecycleActionResult: "CONTINUE", 
            LifecycleActionToken: lifeCycleActionInfo.LifecycleActionToken, 
            LifecycleHookName: lifeCycleActionInfo.LifecycleHookName
        };
        console.log(`Complete ASLCA Param: ${JSON.stringify(params)}`);
        asCient.completeLifecycleAction(params, (err, data) => {
            if (err) reject(err);
            else fulfill(data);
        });
    });
};
exports.handler = eventProcessor;

