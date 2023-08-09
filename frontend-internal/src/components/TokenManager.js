import React, { useState } from "react"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faClipboard } from "@fortawesome/free-regular-svg-icons"
import TokenAdder from "./TokenAdder"
import TokenModal from "./TokenModal"

export default function TokenManager({ macros, tokens }) {
    const [ modal, setModal ] = useState(null)

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
                    <ul>
                        { tokens && tokens.map(token => (
                            <li key="token">
                                <code>
                                    {/* eslint-disable-next-line */}
                                    <a
                                        href=""
                                        onClick={e => {
                                            e.preventDefault()
                                            setModal(token)
                                        }}
                                    >{ token }</a>
                                </code>
                                <span>&nbsp;&nbsp;</span>
                                {/* eslint-disable-next-line */}
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
                {/*<footer>
                    <button className="outline">
                        Delete expired Tokens
                    </button>
                    <button className="outline">
                        Delete all Tokens
                    </button>
                </footer>*/}
            </article>
            <TokenModal 
                open={modal ? true : false}
                closeCallback={() => setModal(null)}
                tokenName={modal}
            />
        </>
    )
}