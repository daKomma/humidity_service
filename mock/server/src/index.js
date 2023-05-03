import express from 'express';

const app = express();

app.listen(3000, () =>
  console.log('Mock Server is listening on port 3000!'),
);

app.get('/data', (req, res) => {
    console.log("Request Time: " + new Date().toISOString())
    console.log("==========================================")
    res.json({
        "hum": 50.5,
        "temp": 25.5
    });
});