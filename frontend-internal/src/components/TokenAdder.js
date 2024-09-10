import React, { useEffect, useState } from "react"

export default function TokenAdder({ serviceCalls }) {
    const [ selectedServiceCall, setSelectedServiceCall ] = useState("")
    const [ numOfUses, setNumOfUses ] = useState(1)
    const [ expiryDate, setExpiryDate ] = useState("")
    const [ expiryTime, setExpiryTime ] = useState("")
    const [ comment, setComment ] = useState("")
    const [ loading, setLoading ] = useState(false)

    useEffect(() => {
        if (serviceCalls)
            setSelectedServiceCall(serviceCalls[0] || "")
    }, [serviceCalls])

    const addToken = event => {
        event.preventDefault()
        setLoading(true)

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
            "service_call_name": selectedServiceCall,
            "expires": dateTimeIso,
            "uses_max": numOfUses,
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
                    setLoading(false)
                    alert(`An error occured while creating the action: ${r.error || "unknown :("}`)
                } else {
                    console.log(r)
                    window.location.reload()
                }
            })
            .catch(console.error)
    }

    return (
        <details>
            <summary role="button" className="secondary">Generate a new token</summary>
            <form>
                <div className="grid">
                    <div>
                        <label>
                            Select action to execute
                            <select
                                value={selectedServiceCall}
                                onChange={e => setSelectedServiceCall(e.target.value)}
                            >
                                { serviceCalls && serviceCalls.map(serviceCall => (
                                    <option
                                        key={serviceCall}
                                        value={serviceCall}
                                    >{ serviceCall } </option>
                                ))}
                            </select>
                        </label>
                    </div>
                    <div>
                        <label>
                            Maximum uses
                            <input 
                                type="number" 
                                value={numOfUses}
                                onChange={e => {
                                    e.preventDefault()
                                    setNumOfUses(parseInt(e.target.value))
                                }}
                            />
                        </label>
                    </div>
                </div>

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

                <label>
                    Comment (optional, visible for token submitter)
                    <input 
                        type="text" 
                        placeholder="E.g. 'Door token for Svenja'"
                        value={comment}
                        onChange={event => setComment(event.target.value)}
                    />
                </label>

                <p><i>The token value will be automatically generated</i></p>
                
                <button 
                    type="submit"
                    onClick={addToken}
                    aria-busy={ loading ? true : false }
                >Generate</button>
            </form>
        </details>
    )
}