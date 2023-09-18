import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { Logger } from '@nestjs/common';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import * as packageJson from 'package.json';

async function bootstrap() {
  const logger = new Logger();
  const app = await NestFactory.create(AppModule);

  if (process.env.SWAGGER_DOCS_PATH) {
    const config = new DocumentBuilder()
      .setTitle(packageJson.name + ' - ' + process.env.STAGE)
      .setDescription('API Documentation')
      .setVersion(packageJson.version)
      .addBearerAuth()
      .build();
    const document = SwaggerModule.createDocument(app, config);
    SwaggerModule.setup(process.env.SWAGGER_DOCS_PATH, app, document);
  }

  app.enableCors({ origin: true, exposedHeaders: ['*'] });

  const port = process.env.PORT || 3000;

  await app.listen(port);

  logger.log(`listening on port ${port}`);
}
bootstrap();
