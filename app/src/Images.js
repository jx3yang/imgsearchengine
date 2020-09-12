import React from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faTimesCircle } from '@fortawesome/free-solid-svg-icons'

function Images(props) {
    return (
        <div className='fadein'>
            {   props.hasDelete &&
                <div 
                    onClick={() => props.removeImages()} 
                    className='delete'
                >
                    <FontAwesomeIcon icon={faTimesCircle} size='2x' />
                </div>
            }
            {
                props.images.map((image, i) => 
                    <img 
                        src={image.path} 
                        alt='' 
                        onError={() => props.onError()}
                    />
                )
            }
            
        </div>
    )
}

export default Images
