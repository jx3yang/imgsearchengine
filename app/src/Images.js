import React from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTimesCircle } from '@fortawesome/free-solid-svg-icons'

function Images(props) {
    return props.images.map((image, i) =>
        <div key={i} className='fadein'>
            <div 
                onClick={() => props.removeImages()} 
                className='delete'
            >
                <FontAwesomeIcon icon={faTimesCircle} size='2x' />
            </div>
            <img 
                src={image.path} 
                alt='' 
                onError={() => props.onError()}
            />
        </div>
    )
}

export default Images
