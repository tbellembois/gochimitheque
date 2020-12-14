import QrScanner from '/js/qr-scanner.min.js';
QrScanner.WORKER_PATH = '/js/qr-scanner-worker.min.js';
//const video = document.getElementById('qr-video');
window.QrScanner = QrScanner;
//qrScanner = new window.QrScanner(video, result => doQR(qrScanner, result));
//qrScanner.start();