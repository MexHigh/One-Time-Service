import React, { useState, useEffect } from "react";
import MacroManager from "./components/MacroManager";
import TokenManager from "./components/TokenManager";
import { CssBaseline, AppBar, Toolbar, Typography, Container, Grid, Card } from "@mui/material";

export default function App() {
    const [ macros, setMacros ] = useState(null)
    const [ tokens, setTokens ] = useState(null)
    
    useEffect(() => {
        fetch("api/internal/tokens/details")
            .then(r => r.json())
            .then(r => {
                setTokens(r.response)
            })
            .catch(console.error)

        fetch("api/internal/macros")
            .then(r => r.json())
            .then(r => {
                setMacros(r.response)
            })
            .catch(console.error)
    }, [])

    return (
        <>
            <CssBaseline />
            <AppBar position="block" sx={{
                marginBottom: "2em"
            }}>
                <Toolbar>
                    <Typography variant="h5">
                        One Time Service
                    </Typography>
                </Toolbar>
            </AppBar>

            <Container>
                <Grid container spacing={3}>
                    <Grid item xs={6}>
                        <Card>
                            Lol
                        </Card>
                    </Grid>
                    <Grid item xs={6}>
                        <Card>
                            Lol
                        </Card>
                    </Grid>

                </Grid>

                <MacroManager
                    macros={macros}
                />
                <TokenManager
                    macros={macros}
                    tokens={tokens}
                />
            </Container>
        </>
    )
}
