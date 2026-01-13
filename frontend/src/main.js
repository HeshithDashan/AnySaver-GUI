window.download = function () {
    let url = document.getElementById("videoUrl").value;
    if (url === "") {
        alert("කරුණාකර ලින්ක් එකක් ඇතුළත් කරන්න!");
        return;
    }
    document.getElementById("status").innerText = "විස්තර ලබාගනිමින්...";
    window.go.main.App.DownloadVideo(url).then((result) => {
        document.getElementById("status").innerText = result;
    }).catch((err) => {
        document.getElementById("status").innerText = "Error: " + err;
    });
};