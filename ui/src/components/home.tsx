import React from "react";
import {Button} from "@mui/material";
import './stylesheets/shared.css'
import axios from "axios";
import {BACKEND_URI, DEFAULT_HEADERS} from "../constants";
import {IState} from "../reducers";
import {useDispatch, useSelector} from "react-redux";
import {bindActionCreators} from "redux";
import {actionCreators} from '../actions'
import {useNavigate} from "react-router-dom";

export function Home() {

  // @ts-ignore
  const state: IState = useSelector(state1 => state1['login'])
  console.log(state)
  const dispatch = useDispatch();
  const {linkTink} = bindActionCreators(actionCreators, dispatch)
  const navigate = useNavigate()

  const linkTinkOnClick = (e: React.MouseEvent) => {
    axios.post(
      `${BACKEND_URI}/link-tink`,
      {
        headers: DEFAULT_HEADERS,
        sessionId: state.sessionId
      }
    ).then(res => {
      window.open(res.data['link']);
      linkTink(true)
      e.preventDefault()
    }).catch(alert)
  }


  return (
    <div className='App-body'>
        <div>
          <div className="button-container">
            <Button type={"submit"} onClick={linkTinkOnClick}>
              Link Tink
            </Button>
          </div>
        </div>
      <div className="button-container">
        <Button type={"submit"}>
          Sync Ledger
        </Button>
      </div>
      <div className="button-container">
        <Button type={"submit"}
                onClick={() => navigate('/classify')}
        >
          Classify
        </Button>
      </div>
      <div className="button-container">
        <Button type={"submit"}>
          Open Fava
        </Button>
      </div>
    </div>
  );
}

export default Home
