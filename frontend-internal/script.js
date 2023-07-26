function test() {
    alert("Lol")
}

window.addEventListener("DOMContentLoaded", () => {
    fetch("api/internal/ping")
        .then(r => r.text())
        .then(r => {
            document.querySelector("#test").innerHTML = r
        })
        .catch(console.error)
})