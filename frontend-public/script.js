//const DEBUG_API_HOST = null
const DEBUG_API_HOST = "http://localhost:1337/"

window.addEventListener("DOMContentLoaded", () => {
    // get DOM elements
    const loadingCard = document.querySelector("#loading-card")
    const enterTokenCard = document.querySelector("#enter-token-card")
    const tokenInvalidCard = document.querySelector("#token-invalid-card")
    const detailsCard = document.querySelector("#details-card")
    const resultCard = document.querySelector("#result-card")

    // get token param
    let params = new URLSearchParams(window.location.search)
    let tokenParam = params.get("token") // null = not set; empty string = set, but empty

    // add event listeners to buttons
    enterTokenCard.querySelector("footer > button").addEventListener("click", () => {
        let tokenVal = enterTokenCard.querySelector("div > input").value
        window.location.assign(`?token=${tokenVal}`)
    })
    let reload = function() { window.location.assign("/") }
    tokenInvalidCard.querySelector("footer > button").addEventListener("click", reload)
    resultCard.querySelector("footer > button").addEventListener("click", reload)
    detailsCard.querySelector("footer > button").addEventListener("click", () => {
        fetch(`${DEBUG_API_HOST ? DEBUG_API_HOST : ""}api/public/token/submit?token=${tokenParam}`, {
            method: "POST"
        })
            .then(async (r) => {
                detailsCard.hidden = true
                resultCard.hidden = false
                
                let jsonVal = await r.json()
                if (jsonVal.error) {
                    resultCard.querySelector("p").innerHTML = "Error: " + jsonVal.error
                } else {
                    resultCard.querySelector("p").innerHTML = "Success!"
                    setTimeout(() => {
                        window.location.assign("/")    
                    }, 2000)
                }
            })
            .catch(console.error)
    })

    if (tokenParam) {
        fetch(`${DEBUG_API_HOST ? DEBUG_API_HOST : ""}api/public/token/details?token=${tokenParam}`)
            .then(r => {
                if (r.status === 404) {
                    loadingCard.hidden = true
                    tokenInvalidCard.hidden = false
                }
                return r.json()
            })
            .then(r => {
                if (!r.response) return
                if (r.response.error) {
                    throw new Error("response has error: " + r.response.error)
                }
                let { macro_name, expires, comment } = r.response

                document.querySelector("#action-container > code").innerHTML = macro_name

                if (expires) {
                    let container = document.querySelector("#expires-container")
                    container.hidden = false
                    container.querySelector("code").innerHTML = new Date(expires).toLocaleString()
                } 
                
                if (comment) {
                    let container = document.querySelector("#comment-container")
                    container.hidden = false
                    container.querySelector("code").innerHTML = comment
                }
                
                loadingCard.hidden = true
                detailsCard.hidden = false
            })
            .catch(console.error)
    } else {
        loadingCard.hidden = true
        enterTokenCard.hidden = false
    }
})