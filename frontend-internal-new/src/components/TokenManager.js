import React from "react"
import TokenAdder from "./TokenAdder"
import TokenEntry from "./TokenEntry"
import { Card, List } from "@mui/material"

export default function TokenManager({ macros, tokens }) {
    return (
        <>
            <Card>
                <header>
                    Manage <span
                        data-tooltip="One time code that will call a Macro once executed"
                        data-placement="right"
                    >Tokens</span>
                </header>
                <List>
                    <TokenAdder macros={macros} />
                    <div>
                        { tokens && Object.entries(tokens).map(([token, details]) => (
                            <TokenEntry
                                key={token}
                                token={token}
                                details={details}
                            />
                        ))}
                    </div>
                </List>
            </Card>
        </>
    )
}