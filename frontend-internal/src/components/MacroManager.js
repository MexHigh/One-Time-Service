import React, { useEffect, useState } from "react"
import MacroAdder from "./MacroAdder"
import MacroModal from "./MacroModal"

export default function MacroManager({ macros }) {
    const [ modal, setModal ] = useState(null)

    return (
        <>
            <article>
                <header>
                    Manage <span 
                        data-tooltip="Macros map a name to a HA service call"
                        data-placement="right"
                    >Macros</span>
                </header>
                <div>
                    <MacroAdder />
                    <div>
                        { macros && macros.map(macro => (
                            <a
                                key={macro}
                                role="button"
                                href=""
                                className="outline"
                                style={{
                                    marginRight: "20px",
                                    marginBottom: "20px"
                                }}
                                onClick={e => {
                                    e.preventDefault()
                                    setModal(macro)
                                }}
                            >{ macro }</a>
                        ))}
                    </div>
                </div>
            </article>
            <MacroModal
                open={modal ? true : false}
                closeCallback={() => setModal(null)}
                macroName={modal}
            />
        </>
    )
}