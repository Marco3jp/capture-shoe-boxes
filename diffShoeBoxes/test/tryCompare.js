#!/bin/node
const {execSync} = require('child_process');
const {performance} = require('perf_hooks');

// disable PAE, cuz 4 images return 65535
const MetricTypes = ["AE", "DSSIM", "Fuzz", "MAE", "MEPP", "MSE", "NCC", "PHASH", "PSNR", "RMSE", "SSIM"];
const VoidImage = ["./void.jpg"];
const TargetImages = ["./exist_1.jpg", "./exist_2.jpg", "./exist_3.jpg", "./exist_4.jpg"];

const cropParams = {
    path: "/tmp/",
    prefix: "cropped_",
    croppedName: {
        voidImage: "",
        targetImages: [],
    },
    x: 380,
    y: 0,
    width: 520,
    height: 720,
};

console.log("========================");
console.log("[Image Description]");
console.log("exist_1.jpg: 5/20");
console.log("exist_2.jpg: 3/20");
console.log("exist_3.jpg: 7/20");
console.log("exist_4.jpg: 2/20");
console.log("===Starting Crop Images===");
try {
    cropParams.croppedName.voidImage = cropParams.prefix + VoidImage[0].split("/")[1];
    let t = execSync(`convert ${VoidImage[0]} -crop ${cropParams.width}x${cropParams.height}+${cropParams.x}+${cropParams.y} ${cropParams.path}${cropParams.croppedName.voidImage}`, {
        cwd: "./",
        stdio: "pipe"
    });
    TargetImages.forEach((image) => {
        const index = cropParams.croppedName.targetImages.push(cropParams.prefix + image.split("/")[1]);
        t = execSync(`convert ${image} -crop ${cropParams.width}x${cropParams.height}+${cropParams.x}+${cropParams.y} ${cropParams.path}${cropParams.croppedName.targetImages[index - 1]}`, {
            cwd: "./",
            stdio: "pipe"
        });
    })
} catch (e) {
    console.error(e);
} finally {
    console.log("  >>> Finished Crop Images");
}
console.log("===\nOUTPUT SAMPLE... [TARGET] METRIC: SCORE(PROCESS_TIME) -- startCount: START_TIME\n===");
MetricTypes.forEach((metric, index) => {
    console.log(`[${metric}] (${index + 1}/${MetricTypes.length})`);
    let totalTime = 0;
    TargetImages.forEach((image, index) => {
        const startTime = performance.now();
        let stdout;
        try {
            stdout = execSync(`compare -metric ${metric} ${cropParams.path}${cropParams.croppedName.voidImage} ${cropParams.path}${cropParams.croppedName.targetImages[index]} /dev/null`, {
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
