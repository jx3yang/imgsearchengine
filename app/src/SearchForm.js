import React, { Component } from 'react'
import { API_URL } from './config'
import Images from './Images'

const knnQueryType = "KNN"
const rangeQueryType = "Range Search"

class SearchForm extends Component {
    constructor(props) {
        super(props)
        this.sourceImage = props.sourceImage
        this.state = {
            queryType: knnQueryType,
            query: 0,
            queried: false,
            resultImages: []
        }
    }

    onSubmit = e => {
        e.preventDefault()
        if (this.state.queryType !== knnQueryType && this.state.queryType !== rangeQueryType) {
            return
        }

        const formData = new FormData()
        formData.append("query", this.state.query)
        formData.append("image", this.sourceImage.path)

        const fetchImages = apiPath => {
            fetch(apiPath, {
                method: 'POST',
                body: formData
            })
            .then(res => {
                if (!res.ok) {
                    throw res
                }
                return res.json()
            })
            .then(images => {
                this.setState({
                    query: 0,
                    queried: true,
                    resultImages: images.map(image => {return image.imageInfo})
                })
            })
        }

        if (this.state.queryType === knnQueryType) {
            fetchImages(`${API_URL}/knn`)
        }

        if (this.state.queryType === rangeQueryType) {
            fetchImages(`${API_URL}/rangesearch`)
        }
    }

    onInputChange = e => {
        this.setState({query: e.target.value})
    }

    onSelectChange = e => {
        this.setState({queryType: e.target.value})
    }

    removeImages = () => {
        this.setState({images: []})
        this.props.removeImage()
    }

    render() {
        const content = () => {
            if (this.state.queried) {
                return (
                    <Images
                        images={this.state.resultImages}
                        removeImages={() => {}}
                        onError={() => {}}
                        hasDelete={false}
                    />
                )
            } else {
                const step = () => {
                    if (this.state.queryType === knnQueryType) {
                        return "1"
                    }
                    return "any"
                }
                return (
                    <form onSubmit={this.onSubmit} className="fadein">
                        <select value={this.state.queryType} onChange={this.onSelectChange}>
                            <option value={knnQueryType}>KNN</option>
                            <option value={rangeQueryType}>Range Search</option>
                        </select>
                        <input 
                            type="number" 
                            step={step()}
                            value={this.state.query}
                            onChange={this.onInputChange}
                            min="0"
                            required
                        />
                        <input type="submit" value="Search" />
                    </form>
                )
            }
        }
        return (
            <div>
                <Images
                    images={[this.sourceImage]}
                    removeImages={this.removeImages}
                    onError={this.props.onError}
                    hasDelete={true}
                />
                {content()}
            </div>
        )
    }
}

export default SearchForm
