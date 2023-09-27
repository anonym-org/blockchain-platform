import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { BlocksModule } from './modules/blocks/blocks.module';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: process.env.ENV_PATH,
    }),
    BlocksModule,
  ],
  controllers: [AppController],
})
export class AppModule {}
