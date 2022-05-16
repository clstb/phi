import React, {useState} from "react";
import '../styles/styles.css'
import {Box, Button, FilledInput, Typography} from "@mui/material";
import axios from "axios";
import {
  ACCESS_TOKEN,
  CORE_URI,
  DEFAULT_HEADERS,
  FAVA_URI,
  LINK_PATH,
  SYNC_PATH,
  TOKEN_PATH,
  USERNAME
} from "../constants";
import {useInput} from "./util";


export function TokenPage() {

  const [error, setError] = useState<string>();

  if (error) {
    throw Error(error)
  }


  const [codeRequested, setCodeRequested] = useState<boolean>(false)
  const {value: accessCode, bind: bindAccessCode} = useInput()
  const [token, setToken] = useState(sessionStorage.getItem(ACCESS_TOKEN))


  const getLink = () => {
    axios.post(
      `${CORE_URI}${LINK_PATH}`,
      {
        headers: DEFAULT_HEADERS,
        test: true
      })
      .then(
        res => {
          console.log(res.data)
          const link = res.data['link']
          setCodeRequested(true)
          window.open(link)
        }
      )
      .catch(err => {
        setError(err.response.statusText)
      })
  }

  const exchangeCodeForToken = () => {
    axios.post(
      `${CORE_URI}${TOKEN_PATH}`,
      {
        headers: DEFAULT_HEADERS,
        access_code: accessCode
      })
      .then(
        res => {
          console.log(res.data)
          const token = res.data['access_token']
          sessionStorage.setItem(ACCESS_TOKEN, token)
          setToken(token)
        }
      )
      .catch(err => {
        setError(err.response.statusText)
      })

  }

  const openFava = () => {
    const uname = sessionStorage.getItem(USERNAME)
    axios.post(
      `${CORE_URI}${SYNC_PATH}`,
      {
        headers: DEFAULT_HEADERS,
        username : uname,
        access_token: token
      })
      .then(
        res => {
          console.log(res.data)
          window.open(FAVA_URI)
        }
      )
      .catch(setError)
  }

  if (token) {
    return (
      <div className='token-body'>
        {"Your access token is \n"}
        <div className={'token-container'}>
          {token}
        </div>
        <Button type={"submit"} sx={{minWidth: "350px"}} onClick={openFava}>
          OK, take me to FAVA ...
        </Button>
      </div>
    );
  }

  if (codeRequested) {
    return (
      <div className='App-body'>
        <Box>
          <div className="input-container">
            <FilledInput type="text"
                         required={true}
                         placeholder={"Paste your access code here"}
                         {...bindAccessCode}

            />
          </div>
          <Button type={"submit"} onClick={exchangeCodeForToken}>
            OK
          </Button>
        </Box>
      </div>
    );
  }





  return (
    <div className='App-body'>
      <Box>
        <Typography variant="subtitle1" component="div">
          First, we'd need to request access code for you
        </Typography>
        <Button type={"submit"} sx={{minWidth: "350px"}} onClick={getLink}>
          Let's go
        </Button>
      </Box>
    </div>
  );

}

export default TokenPage
