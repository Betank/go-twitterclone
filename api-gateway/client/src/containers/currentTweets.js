import { connect } from 'react-redux'
import TweetList from '../components/tweetList'


const mapStateToProps = (state) => {
    return {
        tweets: state.tweets
    }
}

const CurrentTweetList = connect(
    mapStateToProps
)(TweetList)

export default CurrentTweetList