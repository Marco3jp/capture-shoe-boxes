#!/bin/node
const {exec, execSync} = require('child_process');
const {performance} = require('perf_hooks');

// disable PAE, cuz 4 images return 65535
const MetricTypes = ["AE", "DSSIM", "Fuzz", "MAE", "MEPP", "MSE", "NCC", "PHASH", "PSNR", "RMSE", "SSIM"];
const VoidImage = ["./void.jpg"];
const TargetImages = ["./exist_1.jpg", "./exist_2.jpg", "./exist_3.jpg", "./exist_4.jpg"];

console.log("========================");
console.log("[Image Description]");
console.log("exist_1.jpg: 5/20");
console.log("exist_2.jpg: 3/20");
console.log("exist_3.jpg: 7/20");
console.log("exist_4.jpg: 2/20");
console.log("========================\nOUTPUT SAMPLE... [TARGET] METRIC: SCORE(PROCESS_TIME) -- startCount: START_TIME\n========================");
MetricTypes.forEach((metric) => {
    console.log(`[${metric}]`);
    let totalTime = 0;
    TargetImages.forEach((image) => {
        const startTime = performance.now();
        let stdout;
        try {
            stdout = execSync(`compare -metric ${metric} ${VoidImage} ${image} /dev/null`, {
                cwd: "./",
                stdio: "pipe"
            });
        } catch (e) {
            // wtf 画像の完全一致以外エラーを返すっぽい
            const processTime = performance.now() - startTime;
            console.log(`${image}: ${e.stderr}(${processTime.toFixed(3)}ms) -- startCount: ${startTime.toFixed(3)}`);
            totalTime += processTime;
        } finally {
            if (typeof stdout !== "undefined") {
                console.log(`${metric}: ${stdout}`);
            }
        }
    });
    console.log(`  >>> TotalTime: [${totalTime.toFixed(3)} ms]  AveTime: [${(totalTime / TargetImages.length).toFixed(3)} ms]`);
});
