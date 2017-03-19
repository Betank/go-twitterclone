import React, {Component} from 'react'
import { Panel } from 'react-bootstrap'
import Profile from './profile'

class Info extends Component {
    render() {
        return (
            <Panel>
                <Profile></Profile>
            </Panel>
        )
    }
}

export default Info