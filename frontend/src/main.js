let currentSavePath = "";

window.runtime.EventsOn("download-progress", (percent) => {
    const fill = document.getElementById("progress-fill");
    const percentText = document.getElementById("percent");
    const statusText = document.getElementById("status");

    if (fill) {
        fill.style.width = percent + "%";
    }

    if (percentText) {
        percentText.innerText = percent + "%";
    }

    if (statusText && percent < 100) {
        statusText.innerText = "බාගත වෙමින් පවතී...";
    }
});

window.selectFolder = function () {
    window.go.main.App.SelectFolder().then((path) => {
        if (path) {
            currentSavePath = path;

            document.getElementById("selectedPath").innerText = path;
        }
    }).catch((err) => {
        console.error("Folder selection failed:", err);
    });
};

window.download = function () {
    let url = document.getElementById("videoUrl").value;
    
    if (url === "") {
        alert("කරුණාකර ලින්ක් එකක් ඇතුළත් කරන්න!");
        return;
    }

    const fill = document.getElementById("progress-fill");
    const percentText = document.getElementById("percent");
    const statusText = document.getElementById("status");

    if (fill) fill.style.width = "0%";
    if (percentText) percentText.innerText = "0%";
    if (statusText) statusText.innerText = "විස්තර ලබාගනිමින්...";

    window.go.main.App.DownloadVideo(url, currentSavePath).then((result) => {
        if (statusText) {
            statusText.innerText = result;
        }

        if (fill) fill.style.width = "100%";
        if (percentText) percentText.innerText = "100%";
    }).catch((err) => {
        if (statusText) {
            statusText.innerText = "Error: " + err;
        }
    });
};