import React, {Component, PropTypes} from 'react'
import {connect} from 'react-redux'
import {Media, Row, Col, Grid, Panel, ListGroup, ListGroupItem} from 'react-bootstrap'
import {fetchProfile} from '../actions/twitterActions'
import TweetInput from './tweetInput'
import ProfileStatistics from './profilStatistics'

class Profile extends Component {

    componentDidMount() {
        this.props.dispatch(fetchProfile())
    }

    render() {
        return (
            <Panel>
                <ListGroup fill>
                    <ListGroupItem>
                        <Media>
                            <Media.Body>
                                 <Media.Heading>{this.props.user.name}</Media.Heading>
                            </Media.Body>
                        </Media>
                    </ListGroupItem>
                    <ListGroupItem>
                        <ProfileStatistics/>
                    </ListGroupItem>
                </ListGroup>
                <TweetInput/>
            </Panel>
        )
    }
}

Profile.propTypes = {
    user: PropTypes.shape({
        name: PropTypes.string.isRequired
    }).isRequired
}

const mapStateToProps = (state) => {
    return {
        user: state.user
    }
}

export default connect(mapStateToProps)(Profile)