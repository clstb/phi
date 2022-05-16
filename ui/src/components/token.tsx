import React, {useState} from "react";
import '../styles/styles.css'
import {Box, Button, FilledInput, Input, TextField, Typography} from "@mui/material";
import axios from "axios";
import {CORE_URI, DEFAULT_HEADERS} from "../constants";
import {useInput} from "./login";
import {useNavigate} from "react-router-dom";


export function TokenPage() {

  const [error, setError] = useState();

  const navigate = useNavigate()

  if (error) {
    throw Error(error)
  }


  const [codeRequested, setCodeRequested] = useState(false)
  const {value: accessCode, bind: bindAccessCode} = useInput("")
  const [token, setToken] = useState(sessionStorage.getItem('access_token'))


  const getLink = () => {
    axios.post(
      `${CORE_URI}/auth/link`,
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
      `${CORE_URI}/auth/token`,
      {
        headers: DEFAULT_HEADERS,
        access_code: accessCode
      })
      .then(
        res => {
          console.log(res.data)
          const token = res.data['access_token']
          sessionStorage.setItem('access_token', token)
          setToken(token)
        }
      )
      .catch(err => {
        setError(err.response.statusText)
      })

  }

  if (token) {
    return (
      <div className='token-body'>
        {"Your access token is \n"}
        <div className={'token-container'}>
          {token}
        </div>
        <Button type={"submit"} sx={{minWidth: "350px"}} onClick={() => navigate("/home")}>
          Take me home ...
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
