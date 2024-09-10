import React, { useState } from "react"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faTrashCan } from "@fortawesome/free-regular-svg-icons"
import { faArrowTrendUp } from "@fortawesome/free-solid-svg-icons"

export default function TokenEntry({ token, details }) {
    const [ replenishLoading, setReplenishLoading ] = useState(false)
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

    const replenishTokenUses = token => {
        setReplenishLoading(true)

        if (!window.confirm("Really replenish token uses?")) {
            setReplenishLoading(false)
            return false
        }

        fetch(`api/internal/token/replenish?token=${token}`, {
            method: "POST"
        })
            .then(r => r.json())
            .then(r => {
                setTimeout(() => {
                    window.location.reload()
                }, 1000) // artificial delay ;)
            })
            .catch(err => {
                setReplenishLoading(false)
                console.error(err)
            })
    }

    const isExpired = tokenDetails => {
        if (tokenDetails.uses_left <= 0) 
            return true
        if (tokenDetails.expires && Date.now() > new Date(tokenDetails.expires))
            return true
        return false
    }

    const addPreferedLinebreakBeforeToken = shareUrl => {
        let splat = shareUrl.split("=")
        return <span>{ splat[0] }=<wbr/>{ splat[1] }</span>
    }

    return (                            
        <article key={token} style={{
            padding: "1.5em"
        }}>
            <div style={{
                display: "flex",
                justifyContent: "space-between"
            }}>
                <h6 style={{ wordBreak: "break-word" }}>
                    { details.service_call_name }{ isExpired(details) ? " (EXPIRED!)" : "" }
                    { details.comment && (
                        <i> ("{ details.comment }")</i>
                    )}
                </h6>
                <div style={{
                    display: "flex",
                    justifyItems: "end",
                    alignItems: "start"
                }}>
                    <a 
                        aria-busy={ replenishLoading ? true : false }
                        role="button"
                        className="secondary" 
                        style={{
                            marginLeft: "10px",
                            marginBottom: "5px",
                            padding: "5px 10px",
                            cursor: "pointer"
                        }}
                        href=""
                        onClick={e => {
                            e.preventDefault()
                            replenishTokenUses(token)
                        }}
                    >
                        { !replenishLoading && (
                            <FontAwesomeIcon icon={faArrowTrendUp} width={22} />
                        )}
                    </a>
                    <a 
                        aria-busy={ deleteLoading ? true : false }
                        role="button"
                        className="secondary" 
                        style={{
                            marginLeft: "10px",
                            marginBottom: "5px",
                            padding: "5px 10px",
                            cursor: "pointer"
                        }}
                        href=""
                        onClick={e => {
                            e.preventDefault()
                            deleteToken(token)
                        }}
                    >
                        { !deleteLoading && (
                            <FontAwesomeIcon icon={faTrashCan} width={22} />
                        )}
                    </a>
                </div>
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
                <strong>Uses left: </strong>{ details.uses_left }/{ details.uses_max }<br/>
            </small>
        </article>
    )
}