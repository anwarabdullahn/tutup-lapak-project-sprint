import 'reflect-metadata';
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  const port = parseInt(process.env.PORT || '3002', 10);
  console.log(`Profile service listening on port ${port}`);
  await app.listen(port);
}

bootstrap();

