import React, { useEffect, useState } from "react"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faClipboard } from "@fortawesome/free-regular-svg-icons"

export default function TokenManager() {
    const [ tokens, setTokens ] = useState(null)
    
    useEffect(() => {
        fetch(`api/internal/tokens`)
            .then(r => r.json())
            .then(r => {
                setTokens(r.response)
            })
            .catch(console.error)
    }, [])

    return (
        <article>
            <header>
                Manage <span
                    data-tooltip="One time codes that will call a Macro once executed"
                    data-placement="right"
                >Tokens</span>
            </header>
            <div>
                <h5>Token list</h5>
                <ul>
                    { tokens && tokens.map(token => (
                        <li key="token">
                            <a

                            >{ token }</a>
                            <span>&nbsp;&nbsp;&nbsp;</span>
                            <a
                                className="secondary"
                                style={{
                                    cursor: "pointer"
                                }}
                                onClick={e => {
                                    e.preventDefault()
                                    fetch(`api/internal/token/share-url?token=${token}`)
                                        .then(r => r.json())
                                        .then(r => {
                                            // TODO error checking
                                            navigator.clipboard.writeText(r.response);
                                        })
                                        .catch(console.error)
                                }}
                            >
                                <FontAwesomeIcon icon={faClipboard} />
                            </a>
                        </li>
                    ))}
                </ul>
            </div>
            <footer>
                
            </footer>
        </article>
    )
}