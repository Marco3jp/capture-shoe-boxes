#!/bin/node
const {exec} = require('child_process');
const {performance} = require('perf_hooks');

// disable PAE, cuz 4 images return 65535
const MetricTypes = ["AE", "DSSIM", "Fuzz", "MAE", "MEPP", "MSE", "NCC", "PHASH", "PSNR", "RMSE", "SSIM"];
const TargetImages = ["./exist_1.jpg", "./exist_2.jpg", "./exist_3.jpg", "./exist_4.jpg"];

console.log("========================");
console.log("[Image Description]");
console.log("exist_1.jpg: 5/20");
console.log("exist_2.jpg: 3/20");
console.log("exist_3.jpg: 7/20");
console.log("exist_4.jpg: 2/20");
console.log("========================\nOUTPUT SAMPLE... [TARGET] METRIC: SCORE(PROCESS_TIME) -- startCount: START_TIME\n========================");
MetricTypes.forEach((metric) => {
    TargetImages.forEach((image) => {
        const startTime = performance.now();
        exec(`compare -metric ${metric} ./void.jpg ${image}  /dev/null`, {
            cwd: "./"
        }, (error, stdout, stderr) => {
            if (typeof error !== "undefined") {
                console.log(`[${image}] ${metric}: ${stderr}(${performance.now() - startTime}ms) -- startCount: ${startTime}`); // wtf
            } else if (typeof stderr !== "undefined") {
                console.error(`${metric}[StdErr]: ${stderr}`)
            } else {
                console.log(`${metric}: ${stdout}`)
            }
        })
    })
});
