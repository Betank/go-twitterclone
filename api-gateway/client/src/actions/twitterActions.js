import fetch from 'isomorphic-fetch'

export const addTweet = (text) => {
    return dispatch => {
        dispatch({ type: 'ADD_TWEET' })
        return fetch('/api/tweet/', {
            method: 'POST',
            headers: {
                'Content-Type' : 'application/text; charset=UTF-8'
            },
            body: text
        })
        .then(res => {
            if (res.status != 200) {
                dispatch({ type: 'ADD_TWEET_FAILED' })
            } else {
                dispatch({ type: 'ADD_TWEET_SUCCESS' })
            }
        })
        .catch(err => {
            dispatch({ type: 'ADD_TWEET_FAILED' })
        })
    }
}

export const setVisibilityFilter = (filter) => {
    return {
        type: 'SET_VISIBILITY_FILTER',
        filter
    }
}

export function fetchTweets() {
    return dispatch => {
        return fetch('/api/tweets')
        .then(response => response.json())
        .then(json => dispatch(receiveTweets(json)))
    }
}

function receiveTweets(tweets) {
        return {
        type: 'RECEIVE_TWEETS',
        tweets: tweets
    }
}

export function fetchProfile() {
    return dispatch => {
        return fetch('/api/user')
        .then(response => response.json())
        .then(json => dispatch(receiveProfile(json)))
    }
}

function receiveProfile(user) {
        return {
        type: 'RECEIVE_PROFILE',
        user: user
    }
}

export function fetchStatistics() {
    return dispatch => {
        return fetch('/api/stats/')
        .then(response => response.json())
        .then(json => dispatch(receiveStatistics(json)))
    }
}

function receiveStatistics(stats) {
    return {
        type: 'RECEIVE_STATS',
        stats: stats
    }
}