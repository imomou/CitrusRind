var reporter = require('cucumber-html-reporter');

var defaults = {
    theme: 'bootstrap',
    jsonFile: 'in/report.json',
    output: 'out/report.html',
    launchReport: false,
    metadata: {
        "platform": "Linux"
   }
};

reporter.generate(defaults);