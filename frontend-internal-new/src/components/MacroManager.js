import React, { useState } from "react"
import MacroAdder from "./MacroAdder"
import MacroModal from "./MacroModal"
import { Card, List, Typography } from "@mui/material"

export default function MacroManager({ macros }) {
    const [ modal, setModal ] = useState(null)

    return (
        <>
            <Card sx={{
                padding: "1em"
            }}>
                <Typography variant="h6" gutterBottom>
                    Manage macros
                </Typography>
            </Card>
            <MacroModal
                open={modal ? true : false}
                closeCallback={() => setModal(null)}
                macroName={modal}
            />
        </>
    )
}