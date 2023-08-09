import React, { useState } from "react"

export default function MacroAdder() {
    const [ name, setName ] = useState("")
    const [ yaml, setYaml ] = useState("")
    
    const addMacro = event => {
        event.preventDefault()

        if (name === "" || yaml === "") {
            alert("Name or YAML field is empty")
            return false
        }

        fetch("api/internal/macro", {
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
                    alert(`An error occured white creating the macro: ${r.error || "unknown :("}`)
                } else {
                    console.log(r)
                    window.location.reload()
                }
            })
            .catch(console.error)
    }

    return (
        <details>
            <summary role="button" className="secondary">Add a new Macro</summary>
            <form>
                <label>
                    Macro name
                    <input 
                        type="text" 
                        placeholder="E.g. 'Open front door'"
                        value={name}
                        onChange={event => setName(event.target.value)}
                    />
                </label>

                <label>
                    <span>
                        Service call definition in YAML <i>(<a href="/developer-tools/service" target="_blank">create one here</a>)</i>
                    </span>
                    <textarea 
                        type="text" 
                        rows="7"
                        placeholder="service: homeassistant.restart&#10;data: {}"
                        value={yaml}
                        onChange={event => setYaml(event.target.value)} 
                    />
                </label>
                
                <button 
                    type="submit"
                    onClick={addMacro}
                >Create</button>
            </form>
        </details>
    )
}