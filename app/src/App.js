import React, { Component } from 'react'
import Notifications, { notify } from 'react-notify-toast'
import Spinner from './Spinner'
import SearchForm from './SearchForm'
import Buttons from './Buttons'
import { API_URL } from './config'
import './App.css'

const toastColor = { 
    background: '#505050', 
    text: '#fff' 
}

export default class App extends Component {

    state = {
        loading: true,
        uploading: false,
        image: null
    }

    toast = notify.createShowQueue()

    onChange = e => {
        const errs = []
        const file = Array.from(e.target.files)[0]

        const formData = new FormData()
        const types = ['image/png', 'image/jpeg', 'image/gif']
        
        if (types.every(type => file.type !== type)) {
            errs.push(`'${file.type}' is not a supported format`)
        }

        if (file.size > 150000) {
            errs.push(`'${file.name}' is too large, please pick a smaller file`)
        }

        formData.append("file", file)

        if (errs.length) {
            return errs.forEach(err => this.toast(err, 'custom', 2000, toastColor))
        }

        this.setState({ uploading: true })

        fetch(`${API_URL}/image-upload`, {
            method: 'POST',
            body: formData
        })
        .then(res => {
            if (!res.ok) {
                throw res
            }
            return res.json()
        })
        .then(image => {
            this.setState({
                uploading: false, 
                image: image
            })
        })
        .catch(err => {
            err.json().then(e => {
                this.toast(e.message, 'custom', 2000, toastColor)
                this.setState({ uploading: false })
            })
        })
    }

    removeImage = () => {
        this.setState({ image: null })
    }

    onError = () => {
        this.toast('Oops, something went wrong', 'custom', 2000, toastColor)
        this.removeImage()
    }

    render() {
        const { loading, uploading, image } = this.state

        const content = () => {
            switch(true) {
            case uploading:
                return <Spinner />
            case image !== null:
                return (
                    <SearchForm
                        sourceImage={image}
                        removeImage={this.removeImage}
                        onError={this.onError}
                    />
                )
            default:
                return <Buttons onChange={this.onChange} />
            }
        }

        return (
            <div className='container'>
                <Notifications />
                <div className='buttons'>
                    {content()}
                </div>
            </div>
        )
    }
}
