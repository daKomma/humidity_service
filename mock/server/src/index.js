import express from 'express';

const app = express();

app.listen(process.env.PORT, () =>
  console.log('Mock Server is listening on port '+process.env.PORT+'!'),
);

app.get('/data', (req, res) => {
    console.log("Request Time: " + new Date().toISOString())
    console.log("==========================================")
    res.json({
        "hum": parseFloat(process.env.HUM) || 50.5,
        "temp": parseFloat(process.env.TEMP) || 25.5
    });
});