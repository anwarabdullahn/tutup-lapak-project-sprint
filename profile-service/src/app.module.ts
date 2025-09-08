import { Module, Controller, Get } from '@nestjs/common';
import { ProfileModule } from './profile/profile.module';

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

@Module({
  imports: [ProfileModule],
  controllers: [AppController],
})
export class AppModule {}