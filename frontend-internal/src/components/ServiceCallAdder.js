import React, { useState } from "react"

export default function ServiceCallAdder() {
    const [ name, setName ] = useState("")
    const [ yaml, setYaml ] = useState("")
    const [ loading, setLoading ] = useState(false)
    
    const addServiceCall = event => {
        event.preventDefault()
        setLoading(true)

        if (name === "" || yaml === "") {
            alert("Name or YAML field is empty")
            setLoading(false)
            return false
        }

        fetch("api/internal/service-call", {
            method: "POST",
            headers: {
                "content-type": "application/json"
            },
            body: JSON.stringify({
                "name": name,
                "service_payload_yaml_base64": btoa(yaml)
            })
        })
            .then(r => r.json())
            .then(r => {
                if (r.error) {
                    setLoading(false)
                    alert(`An error occured white creating the service call: ${r.error || "unknown :("}`)
                } else {
                    console.log(r)
                    window.location.reload()
                }
            })
            .catch(console.error)
    }

    return (
        <details>
            <summary role="button" className="secondary">Add a new action</summary>
            <form>
                <label>
                    Action name
                    <input 
                        type="text" 
                        placeholder="E.g. 'Open front door'"
                        value={name}
                        onChange={event => setName(event.target.value)}
                    />
                </label>

                <label>
                    <span>
                        Action definition in YAML <i>(<a href="/developer-tools/action" target="_blank">create one here</a>)</i>
                    </span>
                    <textarea 
                        type="text" 
                        rows="7"
                        placeholder="action: homeassistant.restart&#10;data: {}"
                        value={yaml}
                        onChange={event => setYaml(event.target.value)} 
                    />
                </label>
                
                <button 
                    type="submit"
                    onClick={addServiceCall}
                    aria-busy={ loading ? true : false}
                >Create</button>
            </form>
        </details>
    )
}