function test() {
    alert("Lol")
}

window.addEventListener("DOMContentLoaded", () => {
    fetch("/api/public/ping")
        .then(r => r.text())
        .then(r => {
            document.querySelector("#test").innerHTML = r
        })
        .catch(console.error)
})