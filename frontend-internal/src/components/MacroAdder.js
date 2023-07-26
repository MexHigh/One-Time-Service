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

        fetch("/api/internal/macro", {
            method: "POST",
            headers: {
                "content-type": "application/json"
            },
            body: JSON.stringify({
                "name": name,
                "service_payload_yaml_base64": btoa(yaml)
            })
        })
            .then(r => r.json()) // TODO handle errors
            .then(r => {
                console.log(r)
                window.location.reload()
            })
            .catch(console.error)
    }

    return (
        <details>
            <summary role="button" className="secondary">Add a Macro</summary>
            <form>
                <input 
                    type="text" 
                    placeholder="Name"
                    value={name}
                    onChange={event => setName(event.target.value)}
                />

                <a href="/developer-tools/service">Create YAML definition here</a>
                <textarea 
                    type="text" 
                    rows="7"
                    placeholder="YAML definition for service call"
                    value={yaml}
                    onChange={event => setYaml(event.target.value)} 
                />
                
                <button 
                    type="submit"
                    onClick={addMacro}
                >Create</button>
            </form>
        </details>
    )
}