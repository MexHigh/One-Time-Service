import React, { useState } from "react"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faTrashCan } from "@fortawesome/free-regular-svg-icons"
import TokenAdder from "./TokenAdder"

export default function TokenManager({ macros, tokens }) {
    const [ deleteLoading, setDeleteLoading ] = useState(false)

    const deleteToken = token => {
        setDeleteLoading(true)

        if (!window.confirm("Really delete token?")) {
            setDeleteLoading(false)
            return false
        }

        fetch(`api/internal/token?token=${token}`, {
            method: "DELETE"
        })
            .then(r => r.json())
            .then(r => {
                setTimeout(() => {
                    window.location.reload()
                }, 1000) // artificial delay ;)
            })
            .catch(err => {
                setDeleteLoading(false)
                console.error(err)
            })
    }

    const addPreferedLinebreakBeforeToken = shareUrl => {
        let splat = shareUrl.split("=")
        return <span>{ splat[0] }=<wbr/>{ splat[1] }</span>
    }

    return (
        <>
            <article>
                <header>
                    Manage <span
                        data-tooltip="One time code that will call a Macro once executed"
                        data-placement="right"
                    >Tokens</span>
                </header>
                <div>
                    <TokenAdder macros={macros} />
                    <div>
                        { tokens && Object.entries(tokens).map(([token, details]) => (
                            <article key={token} style={{
                                padding: "1.5em"
                            }}>
                                <div style={{
                                    display: "flex",
                                    justifyContent: "space-between"
                                }}>
                                    <h6>
                                        { details.macro_name }
                                        { details.comment && (
                                            <i> ("{ details.comment }")</i>
                                        )}
                                    </h6>
                                    <a 
                                        aria-busy={ deleteLoading ? true : false }
                                        role="button"
                                        className="secondary" 
                                        style={{
                                            alignSelf: "start",
                                            padding: "5px 12px",
                                            cursor: "pointer"
                                        }}
                                        href=""
                                        onClick={e => {
                                            e.preventDefault()
                                            deleteToken(token)
                                        }}
                                    >
                                        { !deleteLoading && (
                                            <FontAwesomeIcon icon={faTrashCan} />
                                        )}
                                    </a>
                                </div>
                                <figure style={{
                                    marginBottom: ".1em"
                                }}>
                                    <code>
                                        { addPreferedLinebreakBeforeToken(details.share_url) }
                                    </code>
                                </figure>
                                <p style={{
                                    marginBottom: "1em"
                                }}>
                                    <small>
                                        <a
                                            href={ details.share_url }
                                            target="_blank"
                                        >
                                            Open Link
                                        </a>
                                    </small>
                                </p>
                                <small>
                                    <strong>Created: </strong>{ new Date(details.created).toLocaleString() }<br/>
                                    <strong>Expires: </strong>{ details.expires ? new Date(details.expires).toLocaleString() : "never" }<br/>
                                    <strong>Uses left: </strong>1 / 1<br/>
                                </small>
                            </article>
                        ))}
                    </div>
                </div>
                {/*<footer>
                    <button className="outline">
                        Delete expired Tokens
                    </button>
                    <button className="outline">
                        Delete all Tokens
                    </button>
                </footer>*/}
            </article>
        </>
    )
}