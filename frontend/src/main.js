window.download = function () {
    let url = document.getElementById("videoUrl").value;

    if (url === "") {
        alert("කරුණාකර ලින්ක් එකක් ඇතුළත් කරන්න!");
        return;
    }

    document.getElementById("result").innerText = "පොඩ්ඩක් ඉන්න, වැඩේ කෙරෙනවා...";

    window.go.main.App.DownloadVideo(url).then((result) => {
        document.getElementById("result").innerText = result;
    }).catch((err) => {
        document.getElementById("result").innerText = "Error: " + err;
    });
};