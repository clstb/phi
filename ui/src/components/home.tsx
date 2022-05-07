import React, {useContext, useState} from "react";
import {Box, Button} from "@mui/material";
import axios from "axios";
import {CORE_URI, DEFAULT_HEADERS} from "../constants";
import {useNavigate} from "react-router-dom";
import '../styles/styles.css'
import {AppContext} from "../index";

export function Home() {

  const [error, setError] = useState(undefined);

  if (error) {
    throw Error(error)
  }

  const context = useContext(AppContext)

  console.log(context)
  const navigate = useNavigate()

  const linkTinkOnClick = (e: React.MouseEvent) => {
    axios.post(
      `${CORE_URI}/link-tink`,
      {
        headers: DEFAULT_HEADERS,
        sessionId: context.sessionId
      }
    ).then(res => {
      window.open(res.data['link']);
      e.preventDefault()
    }).catch(err => setError(err))
  }

  const syncLedger = (e: React.MouseEvent) => {
    axios.post(
      `${CORE_URI}/sync`,
      {
        headers: DEFAULT_HEADERS,
        username: context.username,
        sessionId: context.sessionId
      }
    )
      .catch(err => setError(err))

  }


  return (
    <div className='App-body'>
      <Box
        sx={{
          bgcolor: '#21252b',
          boxShadow: 1,
          borderRadius: 2,
          p: 2,
          minWidth: 300,
        }}
      >
        <div className="button-container">
          <Button type={"submit"}
                  onClick={linkTinkOnClick}
                  sx={{minWidth: 300}}
          >
            Link Tink
          </Button>
        </div>
        <div className="button-container">
          <Button type={"submit"}
                  sx={{minWidth: 300}}
                  onClick={syncLedger}
            >
            Sync Ledger
          </Button>
        </div>
        <div className="button-container">
          <Button type={"submit"}
                  sx={{minWidth: 300}}
                  onClick={() => navigate('/classify')}
          >
            Classify
          </Button>
        </div>
        <div className="button-container">
          <Button type={"submit"}
                  sx={{minWidth: 300}}
          >
            Open Fava
          </Button>
        </div>
      </Box>
    </div>
  );
}

export default Home
