import React from "react";
import {Box, Button} from "@mui/material";
import axios from "axios";
import {CORE_URI, DEFAULT_HEADERS} from "../constants";
import {IState} from "../reducers";
import {useSelector} from "react-redux";
import {useNavigate} from "react-router-dom";
import '../styles/styles.css'

export function Home() {

  // @ts-ignore
  const state: IState = useSelector(state1 => state1['login'])
  console.log(state)
  const navigate = useNavigate()

  const linkTinkOnClick = (e: React.MouseEvent) => {
    axios.post(
      `${CORE_URI}/link-tink`,
      {
        headers: DEFAULT_HEADERS,
        sessionId: state.sessionId
      }
    ).then(res => {
      window.open(res.data['link']);
      e.preventDefault()
    }).catch(alert)
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
                  sx={{minWidth: 300}}>
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
