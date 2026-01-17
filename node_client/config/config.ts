import { configDotenv } from "dotenv";

configDotenv({ path: "../go_server/.env" });

const config = {
  port: process.env.PORT,
  otpCleanup: process.env.OTP_CLEANUP_SECONDS,
  otpLifeSpan: process.env.OTP_LIFESPAN_SECONDS,
  workers: process.env.WORKERS,
  target: process.env.TARGET,
};

export default config;
