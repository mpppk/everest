import React, {useCallback, useMemo} from 'react';
import logo from './logo.svg';
import './App.css';
import {useDropzone} from "react-dropzone";

const baseStyle = {
    margin: '0 auto',
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    padding: '20px',
    borderWidth: 2,
    borderRadius: 2,
    borderColor: '#eeeeee',
    borderStyle: 'dashed',
    backgroundColor: '#282C34',
    color: '#bdbdbd',
    outline: 'none',
    transition: 'border .24s ease-in-out',
    width: '90%'
};

const activeStyle = {
    borderColor: '#2196f3'
};

const acceptStyle = {
    borderColor: '#00e676'
};

const rejectStyle = {
    borderColor: '#ff1744'
};

function App() {
    const onDrop = useCallback(acceptedFiles => {
        // Do something with the files
        console.log(acceptedFiles);
    }, []);
    const {
        getRootProps,
        getInputProps,
        isDragActive,
        isDragAccept,
        isDragReject
    } = useDropzone({onDrop})

    const style = useMemo(() => ({
        ...baseStyle,
        ...(isDragActive ? activeStyle : {}),
        ...(isDragAccept ? acceptStyle : {}),
        ...(isDragReject ? rejectStyle : {})
    }), [
        isDragActive,
        isDragAccept,
        isDragReject
    ]);

    return (
        <div className="App">
            <header className="App-header">
                <a href="https://github.com/mpppk/everest">
                    <img src={logo} className="App-logo" alt="logo"/>
                </a>
                <p>
                </p>
                <section className="container">
                    <div {...getRootProps({style})}>
                        <input {...getInputProps()} />
                        <p>You can generate static file server by drag 'n' drop files here, or click to select files</p>
                    </div>
                </section>
            </header>
        </div>
    );
}

export default App;
