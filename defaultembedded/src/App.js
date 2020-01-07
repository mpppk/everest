import React from 'react';
import logo from './logo.png';
import './App.css';
// import {useDropzone} from "react-dropzone";
// import {Zlib} from 'zlibjs/bin/zip.min.js';

// const baseStyle = {
//     margin: '0 auto',
//     flex: 1,
//     display: 'flex',
//     flexDirection: 'column',
//     alignItems: 'center',
//     padding: '20px',
//     borderWidth: 2,
//     borderRadius: 2,
//     borderColor: '#eeeeee',
//     borderStyle: 'dashed',
//     backgroundColor: '#282C34',
//     color: '#bdbdbd',
//     outline: 'none',
//     transition: 'border .24s ease-in-out',
//     width: '90%'
// };

// const activeStyle = {
//     borderColor: '#2196f3'
// };
//
// const acceptStyle = {
//     borderColor: '#00e676'
// };
//
// const rejectStyle = {
//     borderColor: '#ff1744'
// };
//
// function stringToByteArray(str) {
//     const array = new (window.Uint8Array !== void 0 ? Uint8Array : Array)(str.length);
//     let i;
//     let il;
//
//     for (i = 0, il = str.length; i < il; ++i) {
//         array[i] = str.charCodeAt(i) & 0xff;
//     }
//
//     return array;
// }

function App() {
    // const onDrop = useCallback(acceptedFiles => {
    //     const formData = new FormData();
    //     const zip = new Zlib.Zip();
    //
    //     const promises = acceptedFiles.map((file) => {
    //         return new Promise((resolve, reject) => {
    //             const fr = new FileReader();
    //             fr.onload = eve => {
    //                 const ab = fr.result;
    //                 zip.addFile(new Uint8Array(ab), {
    //                     filename: stringToByteArray(file.name)
    //                 });
    //                 resolve();
    //             };
    //             fr.onerror = () => {
    //                 reject(fr.error);
    //             };
    //             fr.readAsArrayBuffer(file);
    //         });
    //     });
    //
    //     Promise.all(promises).then(() => {
    //         const blob = new Blob([zip.compress()], {type: "application/zip"});
    //         formData.append("file", blob, 'files.zip');
    //         const request = new XMLHttpRequest();
    //         request.open("POST", "/api/upload");
    //         request.send(formData);
    //     });
    // }, []);
    // const {
    //     getRootProps,
    //     getInputProps,
    //     isDragActive,
    //     isDragAccept,
    //     isDragReject
    // } = useDropzone({onDrop})

    // const style = useMemo(() => ({
    //     ...baseStyle,
    //     ...(isDragActive ? activeStyle : {}),
    //     ...(isDragAccept ? acceptStyle : {}),
    //     ...(isDragReject ? rejectStyle : {})
    // }), [
    //     isDragActive,
    //     isDragAccept,
    //     isDragReject
    // ]);

    return (
        <div className="App">
            <header className="App-header">
                <a href="https://github.com/mpppk/everest">
                    <img src={logo} className="App-logo" alt="logo"/>
                </a>
                <h2>everest</h2>
                <p>static file server with no dependencies</p>
                {/*<section className="container">*/}
                {/*    <div {...getRootProps({style})}>*/}
                {/*        <input {...getInputProps()} />*/}
                {/*        <p>You can generate static file server by drag 'n' drop files here, or click to select files</p>*/}
                {/*    </div>*/}
                {/*</section>*/}
            </header>
        </div>
    );
}

export default App;
