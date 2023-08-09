import React, { useEffect, useState } from "react"

export default function TokenAdder({ macros }) {
    const [ selectedMacro, setSelectedMacro ] = useState("")
    const [ comment, setComment ] = useState("")
    const [ expiryDate, setExpiryDate ] = useState("")
    const [ expiryTime, setExpiryTime ] = useState("")

    useEffect(() => {
        if (macros)
            setSelectedMacro(macros[0] || "")
    }, [macros])

    const addToken = event => {
        event.preventDefault()

        let dateTimeIso = undefined
        if (expiryDate || expiryTime) {
            if (!(expiryDate && expiryTime)) {
                alert("Both expiry date and time must be set!")
                return false
            }
            dateTimeIso = new Date(`${expiryDate}T${expiryTime}`).toISOString()
        }

        let parsedComment = comment === "" ? undefined : comment

        let body = {
            "macro_name": selectedMacro,
            "expires": dateTimeIso,
            "comment": parsedComment
        }
        
        fetch("api/internal/token", {
            method: "POST",
            headers: {
                "content-type": "application/json"
            },
            body: JSON.stringify(body)
        })
            .then(r => r.json()) // TODO handle errors
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
            <summary role="button" className="secondary">Generate a new Token</summary>
            <form>
                <label>
                    Select Macro to execute
                    <select>
                        { macros && macros.map((macro, index) => (
                            <option
                                selected={index === 0 ? true : false}
                                key={macro}
                                value={selectedMacro}
                                onChange={e => {
                                    e.preventDefault()
                                    setSelectedMacro(macro)
                                }}
                            >{ macro } </option>
                        ))}
                    </select>
                </label>

                <label>
                    Comment (optional, visible for token submitter)
                    <input 
                        type="text" 
                        placeholder="E.g. 'Door token for Svenja'"
                        value={comment}
                        onChange={event => setComment(event.target.value)}
                    />
                </label>

                <div className="grid">
                    <div>
                        <label>
                            Expiry date (optional)
                            <input 
                                type="date"
                                value={expiryDate}
                                onChange={e => {
                                    e.preventDefault()
                                    setExpiryDate(e.target.value)
                                }}
                            />
                        </label>
                    </div>
                    <div>
                        <label>
                            Expiry time (optional)
                            <input 
                                type="time"
                                value={expiryTime}
                                onChange={e => {
                                    e.preventDefault()
                                    setExpiryTime(e.target.value)
                                }}
                            />
                        </label>
                    </div>
                </div>

                <p><i>A token will be automatically generated</i></p>
                
                <button 
                    type="submit"
                    onClick={addToken}
                >Generate</button>
            </form>
        </details>
    )
}