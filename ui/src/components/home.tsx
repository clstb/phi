import React, {useContext, useState} from "react";
import {Box, Button} from "@mui/material";
import {useNavigate} from "react-router-dom";
import '../styles/styles.css'
import {AppContext} from "../index";
import axios from "axios";
import {CORE_URI, DEFAULT_HEADERS, FAVA_URI} from "../constants";


const username = sessionStorage.getItem("username")!!
const token = sessionStorage.getItem("access_token")!!

export function Home() {

  const [error, setError] = useState(undefined);
  const [synced, setSynced] = useState(false)

  if (error) {
    throw Error(error)
  }

  const context = useContext(AppContext)

  console.log(context)
  const navigate = useNavigate()

  const sync = () => {
    axios.post(
      `${CORE_URI}/sync-ledger`,
      {
        headers: DEFAULT_HEADERS,
        username: username,
        access_token: token
      })
      .then(
        res => {
          console.log(res.data)
          setSynced(true)
          alert("sync OK")
        }
      )
      .catch(setError)
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
                  sx={{minWidth: 300}}
                  onClick={sync}
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
                  onClick={() => window.open(FAVA_URI)}
          >
            Open Fava
          </Button>
        </div>
      </Box>
    </div>
  );
}

export default Home
