import React, { useEffect, useState } from "react"

export default function ServiceCallModal({ open, closeCallback, serviceCallName }) {
    const [ serviceCall, setServiceCall ] = useState(null)
    const [ deleteLoading, setDeleteLoading ] = useState(false)

    useEffect(() => {
        if (!serviceCallName) return
        fetch(`api/internal/service-call/details?name=${serviceCallName}`)
            .then(r => r.json())
            .then(r => {
                setServiceCall(r.response)
            })
            .catch(console.error)
    }, [serviceCallName])

    const deleteServiceCall = () => {
        setDeleteLoading(true)

        if (!window.confirm("Really delete this action? This will also delete all tokens associated with it!")) {
            setDeleteLoading(false)
            return false
        }

        fetch(`api/internal/service-call?name=${serviceCallName}`, {
            method: "DELETE"
        })
            .then(r => r.json())
            .then(r => {
                setTimeout(() => {
                    window.location.reload()
                }, 750) // artificial delay ;)
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
                            setServiceCall(null)
                        }}
                    ></a>
                    <span>Action details: <strong>{ serviceCallName }</strong></span>
                </header>
                { serviceCall ? (
                    <>
                        <div>
                            <p>This action will execute the following Home Assistant action</p>
                            <br />
                            <pre><code>
                                { JSON.stringify(serviceCall, null, 4) }
                            </code></pre>
                        </div>
                        <footer>
                            <button
                                aria-busy={deleteLoading}
                                onClick={deleteServiceCall}
                            >Delete action and associated tokens</button>
                        </footer>
                    </>
                ) : <progress />}
            </article>
        </dialog>
    )
}