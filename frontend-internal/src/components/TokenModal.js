import React, { useEffect, useState } from "react"

export default function TokenModal({ open, closeCallback, tokenName }) {
    const [ token, setToken ] = useState(null)
    const [ deleteLoading, setDeleteLoading ] = useState(false)

    useEffect(() => {
        if (!tokenName) return
        fetch(`api/internal/token/details?token=${tokenName}`)
            .then(r => r.json())
            .then(r => {
                setToken(r.response)
            })
            .catch(console.error)
    }, [tokenName])

    const deleteToken = () => {
        setDeleteLoading(true)

        if (!window.confirm("Really delete token?")) {
            setDeleteLoading(false)
            return false
        }

        fetch(`api/internal/token?token=${tokenName}`, {
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

    return (
        <dialog open={open}>
            <article>
                <header>
                    {/* eslint-disable-next-line */}
                    <a 
                        aria-label="Close" 
                        className="close"
                        onClick={e => {
                            e.preventDefault()
                            closeCallback()
                            setToken(null)
                        }}
                    ></a>
                    <span>Token details: <strong>{ tokenName }</strong></span>
                </header>
                { token ? (
                    <>
                        <div>
                            <p>
                                <strong>Executes Macro: </strong>
                                <code>{ token.macro_name }</code>
                            </p>
                            <p>
                                <strong>Comment: </strong>
                                <code>{ token.comment || "none" }</code>
                            </p>
                            <p>
                                <strong>Expiry: </strong>
                                <code>{ token.expires ? new Date(token.expires).toLocaleString() : "never" }</code>
                            </p>
                        </div>
                        <footer>
                            <button
                                onClick={e => {
                                    e.preventDefault()
                                    fetch(`api/internal/token/share-url?token=${tokenName}`)
                                        .then(r => r.json())
                                        .then(r => {
                                            // TODO error checking
                                            navigator.clipboard.writeText(r.response)
                                                .then(() => {
                                                    console.log(`Successfully copied ${r.response}`)
                                                })
                                                .catch(e => {
                                                    console.log(`Error while copying ${r.response}: ${e}`)
                                                })
                                        })
                                        .catch(console.error)
                                }}
                            >Copy Token URL</button>
                            <button
                                aria-busy={deleteLoading}
                                onClick={deleteToken}
                            >Delete Token</button>
                        </footer>
                    </>
                ) : <progress />}
            </article>
        </dialog>
    )
}