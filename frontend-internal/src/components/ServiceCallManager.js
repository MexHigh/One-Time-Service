import React, { useState } from "react"
import ServiceCallAdder from "./ServiceCallAdder"
import ServiceCallModal from "./ServiceCallModal"

export default function ServiceCallManager({ serviceCalls }) {
    const [ modal, setModal ] = useState(null)

    return (
        <>
            <article>
                <header>
                    Manage actions
                </header>
                <div>
                    <ServiceCallAdder />
                    <div>
                        { serviceCalls && serviceCalls.map(serviceCall => (
                            // eslint-disable-next-line
                            <a
                                key={serviceCall}
                                role="button"
                                href=""
                                className="outline"
                                style={{
                                    marginRight: "20px",
                                    marginBottom: "20px"
                                }}
                                onClick={e => {
                                    e.preventDefault()
                                    setModal(serviceCall)
                                }}
                            >{ serviceCall }</a>
                        ))}
                    </div>
                </div>
            </article>
            <ServiceCallModal
                open={modal ? true : false}
                closeCallback={() => setModal(null)}
                serviceCallName={modal}
            />
        </>
    )
}