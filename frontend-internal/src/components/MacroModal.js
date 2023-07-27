import React, { useEffect, useState } from "react"

export default function MacroModal({ open, closeCallback, macroName }) {
    const [ macro, setMacro ] = useState(null)
    const [ deleteLoading, setDeleteLoading ] = useState(false)

    useEffect(() => {
        if (!macroName) return
        fetch(`api/internal/macro/details?name=${macroName}`)
            .then(r => r.json())
            .then(r => {
                setMacro(r.response)
            })
            .catch(console.error)
    }, [macroName])

    const deleteMacro = () => {
        setDeleteLoading(true)

        if (!window.confirm("Really delete?")) {
            setDeleteLoading(false)
            return false
        }

        fetch(`api/internal/macro?name=${macroName}`, {
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
                    <a 
                        aria-label="Close" 
                        className="close"
                        onClick={e => {
                            e.preventDefault()
                            closeCallback()
                            setMacro(null)
                        }}
                    ></a>
                    <span>Macro details: <strong>{ macroName }</strong></span>
                </header>
                { macro ? (
                    <>
                        <div>
                            <p>This macro will execute the following Home Assistant service</p>
                            <br />
                            <pre><code>
                                { JSON.stringify(macro, null, 4) }
                            </code></pre>
                        </div>
                        <footer>
                            <button
                                aria-busy={deleteLoading}
                                onClick={deleteMacro}
                            >Delete</button>
                        </footer>
                    </>
                ) : <progress />}
            </article>
        </dialog>
    )
}