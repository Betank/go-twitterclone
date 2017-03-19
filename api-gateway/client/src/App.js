import React, {Component} from 'react';
import {Navbar, Grid, Row, Col} from 'react-bootstrap'
import TweetList from './components/tweetList'
import Profile from './components/profile'

class App extends Component {
    render() {
        return (
            <div>
                <Navbar inverse collapseOnSelect fixedTop fluid>
                    <Navbar.Header>
                        <Navbar.Brand>
                            <a href="#">GoTwitterClone</a>
                        </Navbar.Brand>
                        <Navbar.Toggle />
                    </Navbar.Header>
                </Navbar>
                <Grid fluid>
                    <Row className="show-grid">
                        <Col sm={3} md={2}>
                            <Profile/>
                        </Col>
                        <Col sm={9} md={10}>
                            <TweetList/>
                        </Col>
                    </Row>
                </Grid>
            </div>
        );
    }
}

export default App;
