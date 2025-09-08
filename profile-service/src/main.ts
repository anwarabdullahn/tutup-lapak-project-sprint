import 'reflect-metadata';
import { NestFactory } from '@nestjs/core';
import { Module, Controller, Get } from '@nestjs/common';

@Controller()
class AppController {
  @Get('healthz')
  health() {
    return 'ok';
  }

  @Get()
  root() {
    return { service: 'profile-service', status: 'running' };
  }
}

@Module({ controllers: [AppController] })
class AppModule {}

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const port = parseInt(process.env.PORT || '3002', 10);
  await app.listen(port);
}

bootstrap();

