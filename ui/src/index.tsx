import React from 'react';
import {Alert, AlertTitle, Button, CssBaseline} from "@mui/material";
import * as ReactDOM from 'react-dom/client';
import {ErrorBoundary} from "react-error-boundary";
import App from "./app";


const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

// @ts-ignore
const ErrorFallback = ({error, resetErrorBoundary}) => {
  return (
    <Alert variant="filled" severity="error"
           action={
             <Button color="inherit" size="large"
                     onClick={resetErrorBoundary}
             >
               RETRY
             </Button>
           }
    >
      <AlertTitle>Error</AlertTitle>
      <strong>{error.message}</strong>
    </Alert>
  )
}


root.render(
  <React.StrictMode>
    <CssBaseline>
      <ErrorBoundary
        FallbackComponent={ErrorFallback}
        onReset={() => {
          // reset the state of your app so the error doesn't happen again
        }}
      >
        <App/>
      </ErrorBoundary>
    </CssBaseline>
  </React.StrictMode>
);

